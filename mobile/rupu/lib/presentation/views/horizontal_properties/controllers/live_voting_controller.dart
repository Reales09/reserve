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
  final isLoading = false.obs;
  final errorMessage = RxnString();

  @override
  void onReady() {
    super.onReady();
    refresh();
  }

  Future<void> refresh() async {
    isLoading.value = true;
    errorMessage.value = null;
    try {
      // TODO: Fetch live voting data
    } catch (_) {
      errorMessage.value = 'No se pudo cargar la información de la votación en vivo.';
    } finally {
      isLoading.value = false;
    }
  }
}
