import 'dart:async';

import 'package:get/get.dart';
import 'package:rupu/domain/entities/live_voting_result.dart';
import 'package:rupu/domain/repositories/horizontal_properties_repository.dart';
import 'package:rupu/domain/infrastructure/repositories/horizontal_properties_repository_impl.dart';

class LiveVotingController extends GetxController {
  final int propertyId;
  final int votingGroupId;
  final int votingId;
  final HorizontalPropertiesRepository repository;

  LiveVotingController({
    required this.propertyId,
    required this.votingGroupId,
    required this.votingId,
    HorizontalPropertiesRepository? repository,
  }) : repository = repository ?? HorizontalPropertiesRepositoryImpl();

  static String tagFor(int propertyId, int votingGroupId, int votingId) =>
      'live-voting-$propertyId-$votingGroupId-$votingId';

  final liveVotingResult = Rxn<LiveVotingResult>();
  final units = <LiveVotingUnit>[].obs;
  final isLoading = false.obs;
  final errorMessage = RxnString();

  @override
  void onReady() {
    super.onReady();
    refresh();
  }

  late StreamSubscription<LiveVotingResult> _subscription;

  Future<void> refresh() async {
    isLoading.value = true;
    errorMessage.value = null;
    _subscription = repository
        .getLiveVotingStream(
      propertyId: propertyId,
      groupId: votingGroupId,
      votingId: votingId,
    )
        .listen((result) {
      liveVotingResult.value = result;
      units.assignAll(result.units);
      isLoading.value = false;
    }, onError: (error) {
      errorMessage.value = 'No se pudo cargar la información de la votación en vivo.';
      isLoading.value = false;
    });
  }

  @override
  void onClose() {
    _subscription.cancel();
    super.onClose();
  }
}
