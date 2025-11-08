import 'package:get/get.dart';
import 'package:get/get.dart';
import 'package:rupu/domain/entities/horizontal_property_voting_groups.dart';
import 'package:rupu/domain/entities/voting.dart';
import 'package:rupu/domain/repositories/horizontal_properties_repository.dart';
import 'package:rupu/domain/infrastructure/repositories/horizontal_properties_repository_impl.dart';

class VotingGroupDetailController extends GetxController {
  final int propertyId;
  final int votingGroupId;
  final HorizontalPropertiesRepository repository;

  VotingGroupDetailController({
    required this.propertyId,
    required this.votingGroupId,
    HorizontalPropertiesRepository? repository,
  }) : repository = repository ?? HorizontalPropertiesRepositoryImpl();

  static String tagFor(int propertyId, int votingGroupId) =>
      'voting-group-detail-$propertyId-$votingGroupId';

  final group = Rxn<HorizontalPropertyVotingGroup>();
  final votings = <Voting>[].obs;
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
      final result = await repository.getVotings(
        propertyId: propertyId,
        groupId: votingGroupId,
      );
      if (result.success) {
        votings.assignAll(result.data);
      } else {
        errorMessage.value = result.message ?? 'No se pudieron cargar las votaciones.';
      }
    } catch (_) {
      errorMessage.value = 'No se pudo cargar la información del grupo de votación.';
    } finally {
      isLoading.value = false;
    }
  }

  Future<void> activateVoting(int votingId) async {
    try {
      await repository.activateVoting(
        propertyId: propertyId,
        groupId: votingGroupId,
        votingId: votingId,
      );
      await refresh();
    } catch (_) {
      // Handle error
    }
  }

  Future<void> deactivateVoting(int votingId) async {
    try {
      await repository.deactivateVoting(
        propertyId: propertyId,
        groupId: votingGroupId,
        votingId: votingId,
      );
      await refresh();
    } catch (_) {
      // Handle error
    }
  }

  Future<bool> updateVoting({
    required int votingId,
    required Map<String, dynamic> data,
  }) async {
    try {
      final result = await repository.updateVoting(
        propertyId: propertyId,
        groupId: votingGroupId,
        votingId: votingId,
        data: data,
      );
      if (result) {
        await refresh();
      }
      return result;
    } catch (_) {
      return false;
    }
  }

  Future<void> deleteVoting(int votingId) async {
    try {
      await repository.deleteVoting(
        propertyId: propertyId,
        groupId: votingGroupId,
        votingId: votingId,
      );
      await refresh();
    } catch (_) {
      // Handle error
    }
  }
}
